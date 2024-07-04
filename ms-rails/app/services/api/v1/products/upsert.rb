module Services
  module Api
    module V1
      module Products
        class Upsert
          attr_accessor :params, :request

          def initialize(params, request)
            @params = params
            @request = request
          end

          def execute
            ActiveRecord::Base.transaction do
              if params[:id].to_i > 0
                update_product
              else
                insert_product
              end
            end

            product
          end

          private

          def product
            @product ||= find_or_initialize_product
          end

          def find_or_initialize_product
            if params[:id].to_i > 0
              Product.find_by(id: params[:id]) || Product.new
            else
              Product.new
            end
          end

          def update_product
            product.name = params[:name] if params[:name].present?
            product.brand = params[:brand] if params[:brand].present?
            product.price = params[:price] if params[:price].present?
            product.description = params[:description] if params[:description].present?
            product.stock = params[:stock] if params[:stock].present?
            product.save!
          end

          def insert_product
            product.id = Product.maximum(:id).to_i + 1
            product.name = params[:name] if params[:name].present?
            product.brand = params[:brand] if params[:brand].present?
            product.price = params[:price] if params[:price].present?
            product.description = params[:description] if params[:description].present?
            product.stock = params[:stock] if params[:stock].present?
            product.save!
          end
        end
      end
    end
  end
end
