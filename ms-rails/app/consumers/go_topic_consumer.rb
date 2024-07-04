# frozen_string_literal: true
class GoTopicConsumer < Karafka::BaseConsumer
  def consume
    messages.each do |message| 
      process_message(message.payload)
    end
  end

  def process_message(payload)
    puts "Raw payload: #{payload.inspect}"
    
    # Check if payload is a string or a hash
    if payload.is_a?(String)
      params = JSON.parse(payload).with_indifferent_access
    elsif payload.is_a?(Hash)
      params = payload.with_indifferent_access
    else
      raise "Unexpected payload format: #{payload.class}"
    end

    request = {}

    upsert_product(params, request)
  rescue JSON::ParserError => e
    puts "Failed to parse JSON: #{e.message}"
  rescue => e
    puts "Failed to process message: #{e.message}"
  end

  def upsert_product(params, request)
    # Check if ID is present in params
    if params[:id].present? && params[:id].to_i > 0
      product = Product.find_by(id: params[:id]) || Product.new
    else
      product = Product.new
    end

    # Update product attributes based on params
    product.name = params[:name] if params[:name].present?
    product.brand = params[:brand] if params[:brand].present?
    product.price = params[:price] if params[:price].present? && params[:price].to_f > 0
    product.description = params[:description] if params[:description].present?
    product.stock = params[:stock] if params[:stock].present? && params[:price].to_i > 0

    # Call for the upsert service
    upsert_service = Services::Api::V1::Products::Upsert.new(product, request)
    product = upsert_service.execute

    puts "Product upserted: #{product.inspect}"
  rescue ActiveRecord::RecordNotFound => e
    puts e.message
  rescue => e
    puts e.message
  end
end